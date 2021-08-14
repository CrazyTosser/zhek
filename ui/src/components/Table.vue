<template>
  <b-row>
    <b-table
      :caption-top="true"
      :fields="fields"
      :items="items"
      hover
      :selectable="main"
      select-mode="single"
      striped
      @row-selected="select"
      selected-variant="primary"
    >
      <template #table-caption
        ><p class="text-center">{{ caption }}</p></template
      >
      <template #head(control)>
        <font-awesome-icon icon="plus" @click="show(-1)" />
      </template>
      <template #cell(control)="data">
        <div class="d-flex justify-content-around">
          <font-awesome-icon icon="edit" @click="show(data.index)" />
          <font-awesome-icon icon="minus" @click="remove(data.item.rn)" />
        </div>
      </template>
    </b-table>
    <b-modal ref="modal" @ok="save">
      <b-form @submit="save">
        <b-form-group label="Название">
          <b-form-input v-model="form.code" />
        </b-form-group>
        <b-form-group label="Параметры">
          <b-input-group>
            <b-form-select v-model="form.select">
              <b-form-select-option v-for="c in cparam" :value="c" :key="c">{{
                c.code
              }}</b-form-select-option>
            </b-form-select>
            <b-form-input v-if="valueble" v-model.number="form.val" />
            <b-input-group-append>
              <b-input-group-text @click="add">
                <font-awesome-icon icon="plus" />
              </b-input-group-text>
            </b-input-group-append>
          </b-input-group>
        </b-form-group>
        <b-form-group v-if="form.params.length > 0">
          <b-list-group>
            <b-list-group-item v-for="(el, i) in form.params" :key="i">
              <b-input-group>
                <b-form-input :value="el.code" disabled />
                <b-form-input
                  v-if="valueble"
                  :value="form.values[i]"
                  disabled
                />
                <b-input-group-append>
                  <b-input-group-text>
                    <font-awesome-icon icon="minus" @click="pop(i)" />
                  </b-input-group-text>
                </b-input-group-append>
              </b-input-group>
            </b-list-group-item>
          </b-list-group>
        </b-form-group>
      </b-form>
    </b-modal>
  </b-row>
</template>

<script>
export default {
  name: "Table",
  props: {
    fields: Array,
    name: String,
    rn: Number,
    caption: String,
    valueble: {
      type: Boolean,
      default: false,
    },
    main: {
      type: Boolean,
      default: false,
    },
    master: String,
  },
  model: {
    prop: "rn",
    event: "select",
  },
  data() {
    return {
      items: [],
      params: [],
      cparam: [],
      etalon: [],
      form: {
        rn: -1,
        select: {},
        val: 0.0,
        code: "",
        params: [],
        values: [],
      },
    };
  },
  methods: {
    async get() {
      let resp = this.main
        ? await fetch("/" + this.name)
        : await fetch("/" + this.name + "?id=" + this.rn);
      this.items = await resp.json();
      if (this.items === null) this.items = [];
    },
    select(row) {
      if (this.main && row.length > 0) this.$emit("select", row[0].rn);
    },
    show(i) {
      if (i > -1) {
        let cur = this.items[i].params;
        let v = [];
        cur.forEach((el) => v.push(el.val));
        this.cparam = [];
        this.params.forEach((el) => {
          if (cur.filter((x) => x.rn === el.rn).length === 0)
            this.cparam.push(el);
        });
        this.$set(this.form, "rn", this.items[i].rn);
        this.$set(this.form, "code", this.items[i].code);
        this.$set(
          this.form,
          "params",
          this.params.filter((x) => !this.cparam.includes(x))
        );
        this.$set(this.form, "values", v);
      } else {
        if (this.master === undefined) {
          this.cparam = this.params;
          this.form = {
            rn: -1,
            prn: this.rn,
            val: 0.0,
            code: "",
            params: [],
            values: [],
            select: {},
          };
        } else {
          let dif = [];
          let v = [];
          this.params.forEach((el) => {
            let d = this.etalon.filter((x) => x.prn === el.rn);
            if (d.length !== 0) {
              dif.push(el);
              v.push(d[0].val);
              console.log(dif);
              console.log(v);
            }
          });
          this.cparam = this.params.filter((x) => !dif.includes(x));
          this.form = {
            rn: -1,
            prn: this.rn,
            val: 0.0,
            code: "",
            params: dif,
            values: v,
            select: [],
          };
        }
      }
      this.$refs.modal.show();
    },
    async remove(i) {
      let modal = await this.$bvModal.msgBoxConfirm(
        "Вы уверены что хотите удалить контроллер?",
        {
          okVariant: "danger",
          okTitle: "Удалить",
          cancelTitle: "Отмена",
        }
      );
      if (modal === true) {
        if (modal === true) {
          let resp = await fetch("/" + this.name, {
            method: "DELETE",
            body: JSON.stringify({ rn: i }),
          });
          if (resp.ok) this.get();
          else console.log(await resp.text());
        }
      }
    },
    async save(e) {
      e.preventDefault();
      if (this.form.params.length < 1 || this.form.code.length < 1) {
        this.$bvToast.toast(`Перечень параметров не должен быть пустой`, {
          title: "Ошибка",
          autoHideDelay: 5000,
        });
        return;
      }
      for (let i = 0; i < this.form.values.length; i++) {
        this.form.params[i].val = this.form.values[i];
      }
      await fetch("/" + this.name, {
        method: this.form.rn === -1 ? "PUT" : "POST",
        body: JSON.stringify(this.form),
      });
      this.get();
    },
    add() {
      this.form.params.push(this.form.select);
      this.cparam.splice(this.cparam.indexOf(this.form.select), 1);
      this.$set(this.form, "select", {});
      if (this.valueble) {
        this.form.values.push(this.form.val);
        this.form.val = 0;
      }
    },
    pop(i) {
      let el = this.form.params.splice(i, 1)[0];
      this.form.values.splice(i, 1);
      this.cparam.push(el);
    },
  },
  async mounted() {
    let resp = await fetch("/param");
    this.params = await resp.json();
    await this.get();
    if (this.master !== undefined) {
      resp = await fetch("/" + this.master + "?id=" + this.rn, {
        method: "OPTIONS",
      });
      this.etalon = await resp.json();
    }
  },
};
</script>

<style scoped></style>
