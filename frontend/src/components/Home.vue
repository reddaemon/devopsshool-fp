<template>
  <div id="datatb">
    <table class="styled-table">
      <thead>
        <tr>
          <th>Date</th>
          <th>Name</th>
          <th>Valute_ID</th>
          <th>NumCode</th>
          <th>CharCode</th>
          <th>Nominal</th>
          <th>Value</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(result, index) in results" v-bind:key="index">
          <td>{{ result[0] }}</td>
          <td>{{ result[1] }}</td>
          <td>{{ result[2] }}</td>
          <td>{{ result[3] }}</td>
          <td>{{ result[4] }}</td>
          <td>{{ result[5] }}</td>
          <td>{{ result[6] }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
<script>
import DataService from "../services/data.service";

export default {
    name: "Home",
    data() {
        return {
            results : [],
        };
    },
    mounted() {
    DataService.getRate().then(
      (response) => {
        this.results = response.data;
      },
      (error) => {
        this.results =
          (error.response &&
            error.response.data &&
            error.response.data.message) ||
          error.message ||
          error.toString();
      }
    );
  },
};
</script>
<style scoped>
.styled-table {
  width: 100%;
  table-layout: fixed;
  overflow: break-word;
  border-collapse: collapse;
  margin: 0px 0;
  font-size: 0.9em;
  font-family: sans-serif;
  min-width: 400px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.15);
}

.styled-table thead tr {
  background-color: #3d3d3d;
  color: #ddd8d8;
  text-align: left;
}

.styled-table th,
.styled-table td {
  padding: 12px 15px;
}

.styled-table tbody tr {
  border-bottom: 1px solid #dddddd;
}

.styled-table tbody tr:nth-of-type(even) {
  background-color: #f3f3f3;
}

.styled-table tbody tr:last-of-type {
  border-bottom: 2px solid #4c4d4c;
}

.datatb {
  position: absolute;
  left: 0px;
  right: 0px;
  top: 0px;
}
</style>
