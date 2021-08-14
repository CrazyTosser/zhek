<template>
  <b-row>
    <b-table :sticky-header="true" :fields="fields" :value="data"/>
  </b-row>
</template>

<script>
export default {
  name: "Statistic",
  data:() => ({
    fields: [
        'param', 'value', 'date'
    ],
    data: [],
    params: []
  }),
  mounted() {

  }
}
</script>

<style scoped>

</style>