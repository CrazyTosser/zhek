<template>
  <div>
    <Table
        v-model="prn"
        caption="Проекты"
        :fields="fields"
        :main="true"
        name="project"
        :valueble="true"
    />
    <Table
        v-if="prn > 0"
        v-model="prn"
        caption="Адреса"
        :fields="fields"
        name="address"
        master="project"
        :valueble="true"
    />
  </div>
</template>

<script>
import Table from "@/components/Table";
export default {
  name: "Project",
  components: { Table },
  data() {
    return {
      prn: 0,
      fields: [
        { key: "control", label: "" },
        { key: "rn", label: "Рег номер" },
        { key: "code", label: "Наименования" },
      ],
    };
  }
};
</script>
