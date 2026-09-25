export interface LoginForm {
    username: string
    password: string
    isSeller: boolean
}

export interface RegisterForm {
    username: string
    password: string
    confirmPassword: string
}

export interface Account {
    id: number
    username: string
    password: string
    createdDate: Date
    isBuyer: boolean
    isSeller: boolean
    totalGame: number
    totalSpend: number
    totalSale: number
}

export interface Game {
    id: number
    releaseDate: Date
    categories: Category[]
    name: string
    description: string
    price: number
    posterUrl: string
    account: Account
}

export interface GameDetail {
    id: number | null
    releaseDate: Date | null
    categories: Category[] | null
    name: string
    description: string
    price: number
    posterUrl: string
    account: Account | null
    hardwareRequirement: Hardware[] | null
    media: Media[] | null
}

export interface Category {
    id: number
    name: string
    isHidden: boolean
}

export interface Hardware {
    id: number
    operatingSystem: string
    processor: string
    memory: string
    graphicsCard: string
    directXVersion: string
    network: string
    storage: string
    soundcard: string
    additionalNote: string
    type: 'min' | 'recommended'
}

export interface Media {
    id: number
    mediaUrl: string
}