const API = "http://localhost:8080";

export function getToken() {
    return localStorage.getItem("token");
}

export function getUserFromToken() {
    const token = getToken();
    if (!token) return null;

    const payload = JSON.parse(atob(token.split(".")[1]));
    return payload;
}

export async function login(email: string, password: string) {
    const res = await fetch(`${API}/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ email, password })
    });

    if (!res.ok) {
        throw new Error("Login failed");
    }

    const data = await res.json();

    if (!data.token) {
        throw new Error("No token received");
    }

    localStorage.setItem("token", data.token);

    return data;
}

export async function getUsers() {
    const token = getToken();

    if (!token) throw new Error("No token");

    const res = await fetch(`${API}/users`, {
        headers: {
            Authorization: `Bearer ${token}`
        }
    });

    if (!res.ok) {
        throw new Error("Unauthorized");
    }

    return res.json();
}

export async function getPhones(userId: number) {
    const res = await fetch(`${API}/phones?user_id=${userId}`, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    });

    return res.json();
}
export async function register(userData: any) {
    const res = await fetch(`${API}/users/create`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${getToken()}`
        },
        body: JSON.stringify(userData)
    });

    const text = await res.text();

    console.log("STATUS:", res.status);
    console.log("RESPONSE:", text);

    if (!res.ok) {
        const text = await res.text();
        console.error("Backend error:", text);
        throw new Error(text);
    }

    return JSON.parse(text);
}