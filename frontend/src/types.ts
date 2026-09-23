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

export interface Category {
    id: number
    name: string
    isHidden: boolean
}