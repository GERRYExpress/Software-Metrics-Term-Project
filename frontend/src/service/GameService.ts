import axios from "axios";

const apiClient = axios.create({
    baseURL: import.meta.env.VITE_BACKEND_URL,
    headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
    }
})

export default {
    getGames() {
        return apiClient.get('/games?_limit=' + 10 + '&_page=' + 1)
    }
}