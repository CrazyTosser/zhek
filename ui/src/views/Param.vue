<template>
  <b-row>
    <b-table striped hover :fields="fields" :items="data" :primary-key="rn" :sort-by="rn">
      <template #head(command)>
        <font-awesome-icon icon="plus" @click="show(-1)"/>
      </template>
      <template #cell(command)="data">
        <div class="d-flex justify-content-around">
          <font-awesome-icon icon="edit" @click="show(data.item.rn)" />
          <font-awesome-icon icon="minus" @click="remove(data.item.rn)"/>
        </div>
      </template>
    </b-table>
    <b-modal ref="form" :ok-only="true" :no-close-on-backdrop="false" ok-title="Сохранить" @ok="onSubmit">
      <b-form @submit="onSubmit">
        <b-form-group label="Наименование">
          <b-form-input v-model="form.code" required/>
        </b-form-group>
        <b-form-group label="Формула">
          <b-textarea v-model="form.formula" rows="3" />
        </b-form-group>
      </b-form>
    </b-modal>
  </b-row>
</template>

<script>
export default {
  name: "Param",
  data() {
    return {
      fields: [
          'command',
        {key: 'rn', label: 'Рег номер', sortable: true},
        {key: 'code', label: 'Наименование', sortable: true},
        {key: 'formula', label: 'Формула расчета'}
      ],
      data: [],
      form: {
        rn: -1,
        code: '',
        formula: ''
      }
    }
  },
  methods: {
    async getData() {
      let resp = await fetch("/param");

      if (resp.ok) {
        this.data = await resp.json();
      }
    },
    async show(rn) {
      if (rn !== -1) {
        this.form.rn = rn;
        let row = this.data.filter(el => el.rn === rn)
        this.form.code = row[0].code;
        this.form.formula = row[0].formula;
      }
      this.$refs.form.show();
    },
    async remove(rn){
      let modal = await this.$bvModal.msgBoxConfirm("Вы уверены что хотите удалить параметр?", {
        okVariant: 'danger',
        okTitle: 'Удалить',
        cancelTitle: 'Отмена'
      });
      if (modal === true) {
        console.log(modal);
        if (modal === true) {
          let resp = await fetch("/param", {
            method: 'DELETE',
            body: JSON.stringify({rn: rn})
          });
          if (resp.ok)
            this.getData();
          else
            console.log(await resp.text());
        }
      }
    },
    async onSubmit(event) {
      event.preventDefault();
      let resp = await fetch("/param", {
        method: this.form.rn === -1 ? 'PUT' : 'POST',
        body: JSON.stringify(this.form)
      });
      if (resp.ok)
        this.getData();
      else
        console.log(await resp.text());
      this.$refs.form.hide();
    }
  },
  mounted() {
    this.getData();
  }
}
</script>

<style scoped>

</style>