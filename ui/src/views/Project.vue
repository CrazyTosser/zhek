<template>
  <div>
    <b-row>
      <b-table :caption-top="true" :fields="afields" :items="addresses" hover striped>
        <template #table-caption><p class="text-center">Адреса</p></template>
        <template #head(control)>
          <font-awesome-icon icon="plus" @click="showAddr(-1)"/>
        </template>
        <template #cell(control)="data">
          <div class="d-flex justify-content-around">
            <font-awesome-icon icon="edit" @click="showAddr(data.index)"/>
            <font-awesome-icon icon="minus" @click="removeAddr(data.item.rn)"/>
          </div>
        </template>
      </b-table>
    </b-row>
    <b-modal ref="proj" @ok="saveProject">
      <b-form @submit="saveProject">
        <b-form-group label="Название">
          <b-form-input v-model="pform.code"/>
        </b-form-group>
        <b-form-group label="Параметры">
          <b-input-group>
            <b-form-select v-model="pform.prn" :options="cparam"/>
            <b-form-input v-model.number="pform.val"/>
            <b-input-group-append>
              <b-input-group-text @click="addp">
                <font-awesome-icon icon="plus"/>
              </b-input-group-text>
            </b-input-group-append>
          </b-input-group>
        </b-form-group>
        <b-form-group v-if="pform.params.length > 0">
          <b-list-group>
            <b-list-group-item v-for="(el, i) in pform.params">
              <b-input-group>
                <b-form-input :value="params.filter(p => p.value === el.prn)[0].text" disabled/>
                <b-form-input :value="el.val" disabled/>
                <b-input-group-append>
                  <b-input-group-text>
                    <font-awesome-icon icon="minus" @click="remp(i, el.prn)"/>
                  </b-input-group-text>
                </b-input-group-append>
              </b-input-group>
            </b-list-group-item>
          </b-list-group>
        </b-form-group>
      </b-form>
    </b-modal>
    <b-modal ref="addr" @ok="saveAddr">
      <b-form @submit="saveAddr">
        <b-form-group label="Адрес">
          <b-form-input v-model="aform.code"/>
        </b-form-group>
        <b-form-group label="Параметры">
          <b-input-group>
            <b-form-select v-model="aform.prn" :options="cparam"/>
            <b-form-input v-model.number="aform.val"/>
            <b-input-group-append>
              <b-input-group-text @click="adda">
                <font-awesome-icon icon="plus"/>
              </b-input-group-text>
            </b-input-group-append>
          </b-input-group>
        </b-form-group>
        <b-form-group v-if="aform.params.length > 0">
          <b-list-group>
            <b-list-group-item v-for="(el, i) in aform.params">
              <b-input-group>
                <b-form-input :value="params.filter(p => p.value === el.prn)[0].text" disabled/>
                <b-form-input :value="el.val" disabled/>
                <b-input-group-append>
                  <b-input-group-text>
                    <font-awesome-icon icon="minus" @click="rema(i, el.prn)"/>
                  </b-input-group-text>
                </b-input-group-append>
              </b-input-group>
            </b-list-group-item>
          </b-list-group>
        </b-form-group>
      </b-form>
    </b-modal>
  </div>
</template>

<script>
export default {
  name: "Project",
  data() {
    return {
      prn: 0,
      params: [],
      cparam: [],
      pparam: [],
      pform: {
        rn: -1,
        prn: -1,
        val: 0.0,
        params: []
      },
      aform: {
        rn: -1,
        prn: -1,
        val: 0.0,
        code: '',
        params: []
      },
      fields: [
        {key: 'control', label: ''},
        {key: 'rn', label: 'Рег номер'},
        {key: 'code', label: 'Наименования'}
      ],
      afields: [
        {key: 'control', label: ''},
        {key: 'rn', label: 'Рег номер'},
        {key: 'code', label: 'Адрес'},
      ],
      projects: [],
      addresses: []
    }
  },
  methods: {
    showAddr(el) {
      if (el > -1) {
        let cur = this.addresses[el].params;
        this.cparam = [];
        // Вытащили все параметры, которых НЕТ в адресных
        for (let p of this.params)
          if (this.pparam.filter(x => p.value === x.prn).length === 0)
            this.cparam.push(p);
        this.$set(this.aform, 'rn', this.addresses[el].rn);
        this.$set(this.aform, 'code', this.addresses[el].code);
        this.$set(this.aform, 'params', this.addresses[el].params);
      } else {
        this.aform = {
          rn: -1,
          prn: this.prn,
          code: '',
          cparam: 0,
          params: this.pparam
        };
      }
      this.$refs.addr.show();
    },
    showProj(el) {
      if (el > -1) {
        let cur = this.projects[el].params;
        this.cparam = [];
        for (const p of this.params) {
          if (cur.filter(c => p.value === c.prn).length === 0)
            this.cparam.push(p);
        }
        this.$set(this.pform, 'rn', this.projects[el].rn);
        this.$set(this.pform, 'code', this.projects[el].code);
        this.$set(this.pform, 'params', this.projects[el].params);
      } else
        this.pform = {
          rn: -1,
          code: '',
          cparam: 0,
          params: []
        };
      this.$refs.proj.show();
    },
    projectSelect(row) {
      this.prn = row[0].rn;
      this.getAddress();
    },
    async getAddress() {
      let resp = await fetch("/address?id=" + this.prn);
      this.addresses = await resp.json();
      resp = await fetch("/project?id="+this.prn, {
        method: 'OPTIONS'
      });
      this.pparam = await resp.json();
    },
    addp() {
      if (this.pform.prn > 0) {
        this.pform.params.push({prn: this.pform.prn, val: this.pform.val});
        this.cparam = this.cparam.filter(el => el.value !== this.pform.prn);
        this.pform.prn = -1;
        this.pform.val = 0.0;
      }
    },
    adda() {
      if (this.aform.prn > 0) {
        this.aform.params.push({prn: this.aform.prn, val: this.aform.val});
        this.cparam = this.cparam.filter(el => el.value !== this.aform.prn);
        this.aform.prn = -1;
        this.aform.val = 0.0;
      }
    },
    remp(i, el) {
      this.pform.params.splice(i, 1);
      this.cparam.push(this.params.filter(p => el === p.value)[0]);
    },
    rema(i, el) {
      this.aform.params.splice(i, 1);
      this.cparam.push(this.params.filter(p => el === p.value)[0]);
    },
    async getProject() {
      let resp = await fetch("/project");
      this.projects = await resp.json();
    },
    async saveProject(event) {
      event.preventDefault();
      let resp = await fetch("/project", {
        method: this.pform.rn === -1 ? 'PUT' : 'POST',
        body: JSON.stringify(this.pform)
      });
      if (resp.ok) {
        this.getProject();
        this.pform = {
          rn: -1,
          prn: -1,
          val: 0.0,
          params: []
        };
      } else
        console.log(await resp.text());
      this.$refs.proj.hide();
    },
    async saveAddr(event) {
      event.preventDefault();
      let resp = await fetch("/address", {
        method: this.aform.rn === -1 ? 'PUT' : 'POST',
        body: JSON.stringify(this.aform)
      });
      if (resp.ok) {
        this.getAddress();
        this.aform = {
          rn: -1,
          prn: -1,
          val: 0.0,
          code: '',
          params: []
        };
      } else
        console.log(await resp.text());
      this.$refs.proj.hide();
    },
    async removeProj(el) {
      let modal = await this.$bvModal.msgBoxConfirm("Вы уверены что хотите удалить контроллер?", {
        okVariant: 'danger',
        okTitle: 'Удалить',
        cancelTitle: 'Отмена'
      });
      if (modal === true) {
        if (modal === true) {
          let resp = await fetch("/project", {
            method: 'DELETE',
            body: JSON.stringify({rn: rn})
          });
          if (resp.ok)
            this.getProject();
          else
            console.log(await resp.text());
        }
      }
    }
  },
  async mounted() {
    let resp = await fetch("/param");
    if (resp.ok) {
      let res = await resp.json();
      for (const re of res)
        this.params.push({value: re.rn, text: re.code});
      this.cparam = this.params;
    }
    this.getProject();
  }
};
</script>