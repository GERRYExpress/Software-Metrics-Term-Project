import apiClient from "./AxiosClient"
import type { GameDetail, GameInput, GameCreateRequest } from "@/types"

export default {
    getGames(perPage: number, page: number) {
        return apiClient.get('/api/v1/games?_limit=' + perPage + '&_page=' + page)
    },
    getGamesByKeyword(perPage: number, page: number, keyword: string) {
        return apiClient.get('/api/v1/games?keyword=' + keyword + '&_limit=' + perPage + '&_page=' + page)
    },
    getGamesByCategory(perPage: number, page: number, category: string) {
        return apiClient.get('/api/v1/games?category=' + category + '&_limit=' + perPage + '&_page=' + page)
    },
    getGamesByKeywordAndCategory(perPage: number, page: number, keyword: string, category: string) {
        return apiClient.get('/api/v1/games?keyword=' + keyword + '&category=' + category + '&_limit=' + perPage + '&_page=' + page)
    },
    getGame(id: number) {
        return apiClient.get('/api/v1/games/' + id)
    },
    getGamesBySeller(userId: number) {
        return apiClient.get('/api/v1/games?sellerId=' + userId)
    },
    createGame(input: GameInput) {
        return apiClient.post<GameDetail>('/api/v1/games', toGameCreateRequest(input))
    },
    updateGame(id: number, input: GameInput) {
        return apiClient.put<GameDetail>('/api/v1/games/' + id, toGameCreateRequest(input))
    },
    deleteGame(id: number, userId: number) {
        return apiClient.delete(`/api/v1/games/${id}?userId=${userId}`)
    }
}

function toGameCreateRequest(input: GameInput): GameCreateRequest {
    return {
        accountId: input.userId,
        releaseDate: input.releaseDate,
        categoryIds: input.categoryIds,
        name: input.name,
        description: input.description,
        price: input.price,
        posterUrl: input.posterUrl,
        hardwareRequirements: input.hardwareRequirements,
        mediaUrls: input.mediaUrls
    }
}