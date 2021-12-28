import { useCookies } from "vue3-cookies";

const { cookies } = useCookies();

class TokenService {
    getLocalRefreshToken() {
        //const user = JSON.parse(localStorage.getItem("user"))
        const user = cookies.get("user")
        return user?.refresh_token;
    }

    getLocalAccessToken() {
        //const user = JSON.parse(localStorage.getItem("user"))
        const user = cookies.get("user")
        return user?.access_token;
    }

    updateLocalAccessToken(token) {
        //let user = JSON.parse(localStorage.getItem("user"))
        let user = cookies.get("user")
        user.access_token = token;
        //localStorage.SetItem("user", JSON.stringify("user"))
        cookies.set("user", user)
    }

    getUser() {
        //return JSON.parse(localStorage.getItem("user"))
        return cookies.get("user")
    }

    setUser(user) {
        //localStorage.setItem("user", JSON.stringify(user));
        let expires = (new Date(Date.now()+ 86400)).toUTCString();
        cookies.set("user", user, expires, "", "")
    }

    removeUser() {
        //localStorage.removeItem("user");
        cookies.remove("user")
    }
}

export default new TokenService();