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
})