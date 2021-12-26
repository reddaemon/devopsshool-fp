import axios from "axios"

const instance = axios.create({
    baseURL: process.env.VUE_APP_ROOT_API,
    headers: {
        "Content-Type": "application/json",
    },
});

export default instance