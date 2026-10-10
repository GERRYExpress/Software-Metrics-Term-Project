import type { GameDetail } from "@/types";
import { defineStore } from "pinia";


export const useGameStore = defineStore('game', {
    state: (): GameDetail => ({
        id: null,
        releaseDate: null,
        categories: null,
        name: '',
        description: '',
        price: 0,
        posterUrl: '',
        account: null,
        hardwareRequirement: null,
        media: null
    }),
    actions: {
        save(gameDetail: GameDetail) {
            this.id = gameDetail.id
            this.releaseDate = gameDetail.releaseDate
            this.categories = gameDetail.categories
            this.name = gameDetail.name
            this.description = gameDetail.description
            this.price = gameDetail.price
            this.posterUrl = gameDetail.posterUrl
            this.account = gameDetail.account
            this.hardwareRequirement = gameDetail.hardwareRequirement
            this.media = gameDetail.media
        },
        reset() {
            this.id = null
            this.releaseDate = null
            this.categories = null
            this.name = ''
            this.description = ''
            this.price = 0
            this.posterUrl = ''
            this.account = null
            this.hardwareRequirement = null
            this.media = null
        }
    }
})