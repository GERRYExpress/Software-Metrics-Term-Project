import apiClient from "./AxiosClient"

export default {
    getGames() {
        return apiClient.get('/api/v1/games?_limit=' + 10 + '&_page=' + 1)
    },
    getGamesByKeyword(keyword: string) {
        return apiClient.get('/api/v1/games?keyword=' + keyword + '&_limit=' + 10 + '&_page=' + 1)
    },
    getGame(id: number) {
        return apiClient.get('/api/v1/games/' + id)
    }
}