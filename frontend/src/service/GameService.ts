import apiClient from "./AxiosClient"

export default {
    getGames(perPage: number, page: number) {
        return apiClient.get('/api/v1/games?_limit=' + perPage + '&_page=' + page)
    },
    getGamesByKeyword(keyword: string) {
        return apiClient.get('/api/v1/games?keyword=' + keyword + '&_limit=' + 10 + '&_page=' + 1)
    },
    getGame(id: number) {
        return apiClient.get('/api/v1/games/' + id)
    }
}