<template>
  <div class="home">
    <Table
        v-model="crn"
        :fields="fields"
        :main="true"
        caption="Контроллеры"
        name="controller"
    />
    <b-row>
      <b-table
          :caption-top="true"
          :fields="dfields"
          :items="items"
          hover
          striped
      >
        <template #table-caption
        ><p class="text-center">Устройства</p></template
        >
        <template #head(control)>
          <font-awesome-icon icon="plus" @click="show(-1)"/>
        </template>
        <template #cell(control)="data">
          <div class="d-flex justify-content-around">
            <font-awesome-icon icon="edit" @click="show(data.index)"/>
            <font-awesome-icon icon="minus" @click="remove(data.item.rn)"/>
          </div>
        </template>
        <template #cell(adr)="data">
          {{data.item.address.text}}
        </template>
      </b-table>
      <b-modal ref="modal" @ok="save">
        <b-form @submit="save">
          <b-form-group label="Идентификатор">
            <b-form-input v-model="form.uid"/>
          </b-form-group>
          <b-form-group label="Адрес размещения">
            <b-form-select v-model="form.address">
              <b-form-select-option v-for="el in addresses" :value="el" >{{el.text}}</b-form-select-option>
            </b-form-select>
          </b-form-group>
          <b-form-group label="Комментарий">
            <b-form-textarea v-model="form.comment" rows="3" max-rows="6" />
          </b-form-group>
        </b-form>
      </b-modal>
    </b-row>
  </div>
</template>

<style scoped>
.home {
  width: 100%;
  height: 90%;
}
</style>

<script>
import Table from "@/components/Table";

export default {
  name: "Controller",
  components: {Table},
  watch: {
    crn() {
      this.get()
    }
  },
  data() {
    return {
      crn: -1,
      fields: [
        {key: "control", label: ""},
        {key: "rn", label: "Рег номер"},
        {key: "code", label: "Наименования"},
      ],
      dfields: [
        {key: "control", label: ""},
        {key: "rn", label: "Рег номер"},
        {key: "uid", label: "Идентификатор"},
        {key: "adr", label: "Адрес размещения"},
        {key: "comment", label: "Коментарий"}
      ],
      items: [],
      form: {
        rn: -1,
        crn: -1,
        uid: '',
        comment: '',
        address: {}
      },
      addresses: []
    };
  },
  methods: {
    show(el) {
      if (el > -1) {
        this.form = this.items[el];
      } else {
        this.form = {
          rn: -1,
          crn: -1,
          uid: '',
          comment: '',
          address: {}
        };
      }
      this.$refs.modal.show();
    },
    async save(e) {
      e.preventDefault();
      let form = this.form;
      form.crn = this.crn;
      form.arn = form.address.value;
      let resp = await fetch('/device', {
        method: form.rn === -1 ? 'PUT' : 'POST',
        body: JSON.stringify(form)
      });
      if (resp.ok) {
        this.get();
        this.$refs.modal.hide();
      }
    },
    async remove(rn) {
      await fetch('/device', {
        method: 'DELETE',
        body: JSON.stringify({rn: rn})
      });
      await this.get();
    },
    async get() {
      let resp = await fetch('/device?id=' + this.crn);
      this.items = await resp.json();
    }
  },
  async mounted() {
    let resp = await fetch("/address", {method: 'OPTIONS'});
    if (resp.ok)
      this.addresses = await resp.json();
    else
      console.error(await resp.text());
    this.get();
  }
};
</script>
