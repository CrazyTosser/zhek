<template>
  <b-container>
    <Table
      caption="События"
      :fields="fields"
      name="event"
      :main="true"
    />
    <b-row>
      <b-form inline>
        <b-form-select v-model="etype"
                       :options="[{value: 1, text: 'Проекты'}, {value: 2, text: 'Здания'}]" />
      </b-form>
      <b-table
        :caption-top="true"
        :fields="lfield"
        :items="events"
        hover
        select-mode="single"
        striped
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
        <template #cell(work)="data">
          <div class="d-flex justify-content-around">
            <font-awesome-icon icon="check" v-if="data.item.work" />
            <font-awesome-icon icon="times" v-if="!data.item.work" />
          </div>
        </template>
    </b-table>
    </b-row>
  </b-container>
</template>

<script>
import Table from "@/components/Table";
export default {
  name: "test",
  components: { Table },
  data() {
    return {
      val: 12,
      etype: 0,
      fields: [
        { key: "control", label: "" },
        { key: "rn", label: "Рег номер" },
        { key: "code", label: "Наименования" },
        { key: "event", label: "Действие"}
      ],
      lfield: [
        { key: "control", label: "" },
        { key: "condition", label: "Условие срабатывания" },
        { key: "work", label: "Включено ли событие"},
      ],
      events: [
        {
          rn: 1,
          condition: "and(<,>)",
          work: true
        }
      ]
    };
  },
};
</script>

<style scoped></style>
