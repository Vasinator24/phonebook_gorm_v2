export interface User {
    id: number;
    username: string;
    email: string;
    names: string;
}

export interface Phone {
    id: number;
    user_id: number;
    number: string;
}