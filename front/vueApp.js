const url = "http://127.0.0.1:8080/v1/rate";

const vm = new Vue({
    el: '#app',
    data: {
        results: []
    },

    computed: {
        sortedItems: function () {
            return this.results.sort((a, b) => new Date(a.date) - new Date(b.date))
        }
    },

    mounted() {
        axios.get(url).then(response => {
            this.results = response.data
            this.results.sort((a, b) => new Date(a.date) - new Date(b.date))
            console.log(this.results)
        })
            .catch(error => console.log(error))
    }
});