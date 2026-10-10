import apiClient from "./AxiosClient"
import type { GameInput } from "@/types"

export default {
    getGames() {
        return apiClient.get('/api/v1/games?_limit=' + 10 + '&_page=' + 1)
    },
    getGamesByKeyword(keyword: string) {
        return apiClient.get('/api/v1/games?keyword=' + keyword + '&_limit=' + 10 + '&_page=' + 1)
    },
    getGame(id: number) {
        return apiClient.get('/api/v1/games/' + id)
    },
    getGamesBySeller(userId: number) {
        return apiClient.get('/api/v1/games?sellerId=' + userId)
    },
    createGame(input: GameInput) {
        return apiClient.post('/api/v1/games', input)
    },
    updateGame(id: number, input: GameInput) {
        return apiClient.put('/api/v1/games/' + id, input)
    },
    deleteGame(id: number, userId: number) {
        return apiClient.delete(`/api/v1/games/${id}?userId=${userId}`)
    }
}