import { useEffect, useState } from "react";
import { getUsers, getPhones, login, register } from "./api/api";
import type { User, Phone } from "./types";

function getUserFromToken() {
    const token = localStorage.getItem("token");
    if (!token) return null;

    try {
        return JSON.parse(atob(token.split(".")[1]));
    } catch {
        return null;
    }
}

function App() {
    const [users, setUsers] = useState<User[]>([]);
    const [phones, setPhones] = useState<Phone[]>([]);
    const [selectedUserId, setSelectedUserId] = useState<number | null>(null);
    const [selectedUserName, setSelectedUserName] = useState("");

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [isLoggedIn, setIsLoggedIn] = useState(!!localStorage.getItem("token"));

    const [page, setPage] = useState(1);
    const USERS_PER_PAGE = 5;

    const [showRegister, setShowRegister] = useState(false);

    const [regUsername, setRegUsername] = useState("");
    const [regEmail, setRegEmail] = useState("");
    const [regNames, setRegNames] = useState("");
    const [regPassword, setRegPassword] = useState("");
    const [phoneInputs, setPhoneInputs] = useState([""]);

    const [editUser, setEditUser] = useState<User | null>(null);
    const [newPhone, setNewPhone] = useState("");

    const user = getUserFromToken();
    const isAdmin = user?.role === "admin";

    useEffect(() => {
        if (isLoggedIn) loadUsers();
    }, [isLoggedIn]);

    async function loadUsers() {
        try {
            const data = await getUsers();
            setUsers(data);
        } catch {
            handleLogout();
        }
    }

    async function loadPhones(userId: number) {
        const data = await getPhones(userId);
        setPhones(data);
    }

    async function handleLogin() {
        try {
            await login(email, password);
            setIsLoggedIn(true);
            setPage(1);
            loadUsers();
        } catch {
            alert("Login failed");
        }
    }

    async function handleRegister() {
        if (!regUsername || !regEmail || !regNames || !regPassword) {
            alert("Fill required fields!");
            return;
        }

        try {
            const filteredPhones = phoneInputs
                .filter(p => p.trim() !== "")
                .map(p => ({ number: p }));

            await register({
                username: regUsername,
                email: regEmail,
                password: regPassword,
                names: regNames,
                phones: filteredPhones
            });

            setShowRegister(false);
            loadUsers();
        } catch {
            alert("Could not register user");
        }
    }

    async function handleDeleteUser(id: number) {
        await fetch(`http://localhost:8080/users?id=${id}`, {
            method: "DELETE",
            headers: {
                Authorization: `Bearer ${localStorage.getItem("token")}`
            }
        });

        loadUsers();
    }

    async function handleUpdateUser() {
        if (!editUser) return;

        await fetch(`http://localhost:8080/users`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${localStorage.getItem("token")}`
            },
            body: JSON.stringify(editUser)
        });

        setEditUser(null);
        loadUsers();
    }

    async function handleAddPhone() {
        if (!selectedUserId || !newPhone) return;

        await fetch(`http://localhost:8080/phones`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${localStorage.getItem("token")}`
            },
            body: JSON.stringify({
                user_id: selectedUserId,
                number: newPhone
            })
        });

        setNewPhone("");
        loadPhones(selectedUserId);
    }

    async function handleDeletePhone(id: number) {
        await fetch(`http://localhost:8080/phones?id=${id}`, {
            method: "DELETE",
            headers: {
                Authorization: `Bearer ${localStorage.getItem("token")}`
            }
        });

        if (selectedUserId) loadPhones(selectedUserId);
    }

    function handleLogout() {
        localStorage.removeItem("token");

        setIsLoggedIn(false);
        setUsers([]);
        setPhones([]);

        setSelectedUserId(null);
        setSelectedUserName("");

        setPage(1);
        setEmail("");
        setPassword("");
    }

    function handleUserClick(u: User) {
        setSelectedUserId(u.id);
        setSelectedUserName(u.names);
        loadPhones(u.id);
    }

    const startIndex = (page - 1) * USERS_PER_PAGE;
    const paginatedUsers = users.slice(startIndex, startIndex + USERS_PER_PAGE);
    const totalPages = Math.ceil(users.length / USERS_PER_PAGE);

    const styles: { [key: string]: React.CSSProperties } = {
        container: {
            padding: 20,
            fontFamily: "Arial",
            background: "#f4f6f8",
            minHeight: "100vh"
        },

        header: {
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            marginBottom: 20,
            padding: 10,
            background: "#fff",
            borderRadius: 10
        },

        grid: {
            display: "grid",
            gridTemplateColumns: "1fr 1fr",
            gap: 20
        },

        card: {
            background: "#fff",
            padding: 20,
            borderRadius: 12,
            boxShadow: "0 4px 12px rgba(0,0,0,0.08)"
        },

        table: {
            width: "100%",
            borderCollapse: "collapse"
        },

        row: {
            cursor: "pointer",
            transition: "0.2s",
        },

        input: {
            flex: 1,
            padding: 10,
            borderRadius: 8,
            border: "1px solid #ddd"
        },

        button: {
            padding: "8px 12px",
            borderRadius: 6,
            border: "none",
            cursor: "pointer"
        },

        addBtn: {
            padding: "8px 12px",
            background: "#2ecc71",
            color: "white",
            border: "none",
            borderRadius: 6,
            cursor: "pointer"
        },

        deleteBtn: {
            padding: "6px 10px",
            background: "#e74c3c",
            color: "white",
            border: "none",
            borderRadius: 6,
            cursor: "pointer"
        },

        editBtn: {
            padding: "6px 10px",
            background: "#3498db",
            color: "white",
            border: "none",
            borderRadius: 6,
            cursor: "pointer",
            marginRight: 5
        },

        dangerBtn: {
            padding: "6px 10px",
            background: "#e74c3c",
            color: "white",
            border: "none",
            borderRadius: 6,
            cursor: "pointer"
        },

        phoneRow: {
            display: "flex",
            justifyContent: "space-between",
            padding: "6px 0"
        },

        pagination: {
            marginTop: 10,
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center"
        },

        modalOverlay: {
            position: "fixed",
            top: 0,
            left: 0,
            width: "100%",
            height: "100%",
            background: "rgba(0,0,0,0.5)",
            display: "flex",
            justifyContent: "center",
            alignItems: "center"
        },

        modal: {
            background: "white",
            padding: 20,
            borderRadius: 12,
            width: 300
        }
    };
    if (!isLoggedIn) {
        return (
            <div style={styles.authContainer}>
                <div style={styles.authCard}>

                    <h1 style={styles.logo}>📱 Phonebook</h1>

                    <p style={styles.subtitle}>
                        Login to continue
                    </p>

                    <input
                        placeholder="Email"
                        value={email}
                        onChange={e => setEmail(e.target.value)}
                        style={styles.authInput}
                    />

                    <input
                        type="password"
                        placeholder="Password"
                        value={password}
                        onChange={e => setPassword(e.target.value)}
                        style={styles.authInput}
                    />

                    <button
                        onClick={handleLogin}
                        style={styles.loginBtn}
                    >
                        Login
                    </button>

                    <div style={{ marginTop: 20 }}>
                        <button
                            onClick={() => setShowRegister(!showRegister)}
                            style={styles.switchBtn}
                        >
                            {showRegister
                                ? "Close Registration"
                                : "Create Account"}
                        </button>
                    </div>

                    {showRegister && (
                        <div style={{ marginTop: 25 }}>

                            <input
                                placeholder="Username"
                                value={regUsername}
                                onChange={e => setRegUsername(e.target.value)}
                                style={styles.authInput}
                            />

                            <input
                                placeholder="Email"
                                value={regEmail}
                                onChange={e => setRegEmail(e.target.value)}
                                style={styles.authInput}
                            />

                            <input
                                type="password"
                                placeholder="Password"
                                value={regPassword}
                                onChange={e => setRegPassword(e.target.value)}
                                style={styles.authInput}
                            />

                            <input
                                placeholder="Full Name"
                                value={regNames}
                                onChange={e => setRegNames(e.target.value)}
                                style={styles.authInput}
                            />

                            {phoneInputs.map((p, i) => (
                                <input
                                    key={i}
                                    placeholder="Phone Number"
                                    value={p}
                                    onChange={e => {
                                        const copy = [...phoneInputs];
                                        copy[i] = e.target.value;
                                        setPhoneInputs(copy);
                                    }}
                                    style={styles.authInput}
                                />
                            ))}

                            <button
                                onClick={() =>
                                    setPhoneInputs([...phoneInputs, ""])
                                }
                                style={styles.addPhoneBtn}
                            >
                                + Add Phone
                            </button>

                            <button
                                onClick={handleRegister}
                                style={styles.registerBtn}
                            >
                                Create Account
                            </button>
                        </div>
                    )}
                </div>
            </div>
        );
    }

    return (
    <div style={styles.container}>
        {/* HEADER */}
        <div style={styles.header}>
            <h2 style={{ margin: 0 }}>📱 Phonebook System</h2>

            <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
                <span style={{ fontWeight: "bold" }}>
                    {user?.username}
                </span>

                <button onClick={handleLogout} style={styles.dangerBtn}>
                    Logout
                </button>
            </div>
        </div>

        <div style={styles.grid}>
            {/* USERS */}
            <div style={styles.card}>
                <h3>Users</h3>

                <table style={styles.table}>
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Username</th>
                            <th>Name</th>
                            {isAdmin && <th>Actions</th>}
                        </tr>
                    </thead>

                    <tbody>
                        {paginatedUsers.map(u => (
                            <tr
                                key={u.id}
                                style={styles.row}
                            >
                                <td onClick={() => handleUserClick(u)}>{u.id}</td>
                                <td onClick={() => handleUserClick(u)}>{u.username}</td>
                                <td onClick={() => handleUserClick(u)}>{u.names}</td>

                                {isAdmin && (
                                    <td>
                                        <button
                                            style={styles.editBtn}
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                setEditUser(u);
                                            }}
                                        >
                                            Edit
                                        </button>

                                        <button
                                            style={styles.deleteBtn}
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                handleDeleteUser(u.id);
                                            }}
                                        >
                                            Delete
                                        </button>
                                    </td>
                                )}
                            </tr>
                        ))}
                    </tbody>
                </table>

                {/* PAGINATION */}
                <div style={styles.pagination}>
                    <button
                        onClick={() => setPage(p => Math.max(1, p - 1))}
                        style={styles.button}
                    >
                        Prev
                    </button>

                    <span>Page {page} / {totalPages}</span>

                    <button
                        onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                        style={styles.button}
                    >
                        Next
                    </button>
                </div>
            </div>

            {/* PHONES */}
            <div style={styles.card}>
                <h3>Phones</h3>

                <div style={{ marginBottom: 10, fontWeight: "bold" }}>
                    {selectedUserName || "Select user"}
                </div>

                {phones.map(p => (
                    <div key={p.id} style={styles.phoneRow}>
                        <span>{p.number}</span>

                        <button
                            style={styles.deleteBtn}
                            onClick={() => handleDeletePhone(p.id)}
                        >
                            X
                        </button>
                    </div>
                ))}

                <div style={{ marginTop: 15, display: "flex", gap: 10 }}>
                    <input
                        placeholder="new phone"
                        value={newPhone}
                        onChange={e => setNewPhone(e.target.value)}
                        style={styles.input}
                    />

                    <button onClick={handleAddPhone} style={styles.addBtn}>
                        Add
                    </button>
                </div>
            </div>
        </div>

        {/* EDIT MODAL */}
        {editUser && (
            <div style={styles.modalOverlay}>
                <div style={styles.modal}>
                    <h3>Edit User</h3>

                    <input
                        value={editUser.username}
                        onChange={e =>
                            setEditUser({ ...editUser, username: e.target.value })
                        }
                        style={styles.input}
                    />

                    <input
                        value={editUser.names}
                        onChange={e =>
                            setEditUser({ ...editUser, names: e.target.value })
                        }
                        style={styles.input}
                    />

                    <input
                        value={editUser.email}
                        onChange={e =>
                            setEditUser({ ...editUser, email: e.target.value })
                        }
                        style={styles.input}
                    />

                    <div style={{ display: "flex", gap: 10 }}>
                        <button onClick={handleUpdateUser} style={styles.addBtn}>
                            Save
                        </button>

                        <button
                            onClick={() => setEditUser(null)}
                            style={styles.deleteBtn}
                        >
                            Cancel
                        </button>
                    </div>
                </div>
            </div>
        )}
    </div>
)};

export default App;