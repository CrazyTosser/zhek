package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"math"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	control = make(map[int][]Param)
)

func initScanner(pool *pgxpool.Pool, ctx context.Context, config *Config) {
	for ctx.Err() == nil {
		// filling controller map
		query, err := pool.Query(ctx, "select c.rn, p.rn, p.code, p.formula, rank() over (partition by c.rn order by cp.pos) "+
			"from controller c inner join controller_parameter cp on c.rn = cp.crn inner join parameter p on p.rn = cp.prn ")
		if err != nil {
			return
		}
		for query.Next() {
			p := Param{}
			i := 0
			err := query.Scan(&i, &p.Rn, &p.Code, &p.Formula)
			if err != nil {
				return
			}
			control[i] = append(control[i], p)
		}
		//fill device map
		query, err = pool.Query(ctx, "select d.rn, d.uid, d.crn, l.rn from devices d "+
			"inner join location l on d.rn = l.drn")
		if err != nil {
			return
		}
		for query.Next() {
			dev := Device{}
			err := query.Scan(&dev.Rn, &dev.Uid, &dev.Crn, &dev.Address)
			if err != nil {
				return
			}
			dev.Schema = control[dev.Crn]
			config.Devices[dev.Uid] = dev
		}
		time.Sleep(time.Minute * 10)
	}
}

func ServiceUdp(ctx context.Context, address string, port string, pool *pgxpool.Pool, config *Config) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", address, port))
	pc, err := net.ListenUDP("udp", addr)
	if err != nil {
		return
	}
	defer func(pc net.PacketConn) {
		_ = pc.Close()
	}(pc)

	for ctx.Err() == nil {
		buffer := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			return
		}
		fmt.Printf("packet-received: bytes=%d from=%s\n",
			n, addr.String())
		go parseAndSave(ctx, pool, buffer, config)
	}
}

func parseAndSave(ctx context.Context, pool *pgxpool.Pool, buf []byte, conf *Config) {
	var uid [6]byte
	var tmp float32
	rd := bytes.NewReader(buf)
	readString, err := rd.Read(buf)
	if err != nil || readString < 6 {
		return
	}
	if d, ex := conf.Devices[string(uid[:])]; ex {
		lrn := d.Address
		// Process packet from device
		var stack []Param //stack for calculating parameters
		for _, param := range d.Schema {
			if len(param.Formula) == 0 {
				tst := make([]byte, 4)
				readString, err := rd.Read(tst)
				if err != nil || readString < 4 {
					return
				}
				bits := binary.LittleEndian.Uint32(tst)
				tmp = math.Float32frombits(bits)
				_, err = pool.Query(ctx, "insert into statistic (lrn, prn, val) values (?, ?, ?)", lrn, param.Rn, tmp)
				if err != nil {
					return
				}
			} else {
				stack = append(stack, param)
			}
		}
		for _, s := range stack {
			pos := strings.Index(s.Formula, "mid(")
			if pos == 0 {
				parse := strings.Split(strings.Trim(s.Formula, "mid()"), ",")
				prn, _ := strconv.Atoi(strings.TrimSpace(parse[0]))
				cnt, _ := strconv.Atoi(strings.TrimSpace(parse[1]))
				q, err := pool.Query(ctx, "select s.val from statistic s "+
					"inner join location l on l.rn = s.lrn "+
					"where l.drn = ? and s.prn = ? "+
					"order by l.date desc limit ?", d.Rn, prn, cnt)
				if err != nil {
					return
				}
				sum := 0.0
				c := 0
				for q.Next() {
					tmp := 0.0
					err := q.Scan(&tmp)
					if err != nil {
						return
					}
					sum += tmp
					c += 1
				}
				_, err = pool.Query(ctx, "insert into statistic (lrn, prn, val) values (?, ?, ?)",
					lrn, s.Rn, sum/float64(c))
				if err != nil {
					return
				}
			}
		}
	}
}
