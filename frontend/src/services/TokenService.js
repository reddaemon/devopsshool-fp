import { useCookies } from "vue3-cookies";

const { cookies } = useCookies();

class TokenService {
    getLocalRefreshToken() {
        const user = cookies.get("user")
        return user?.refresh_token;
    }

    getLocalAccessToken() {
        const user = cookies.get("user")
        return user?.access_token;
    }

    updateLocalAccessToken(token) {
        let user = cookies.get("user")
        user.access_token = token;
        cookies.set("user", user)
    }

    getUser() {
        return cookies.get("user")
    }

    setUser(user) {
        let expires = (new Date(Date.now()+ 86400)).toUTCString();
        cookies.set("user", user, expires, "", "")
    }

    removeUser() {
        cookies.remove("user")
    }
}

export default new TokenService();