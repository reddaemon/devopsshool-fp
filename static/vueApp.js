const url = "http://127.0.0.1:8080/v1/rate";

const vm = new Vue({
    el: '#app',
    data: {
        results: []
    },

    
    mounted() {
        axios.get(url).then(response => {
            this.results = response.data
            console.log(this.results)
        })
            .catch(error => console.log(error))
    }
});