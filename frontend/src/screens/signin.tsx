import "../styles.css";
import PasswordScreen from "../components/password";
import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

const SignInScreen = () => {
	const [error, setError] = useState<string>("");
	const navigate = useNavigate();

	function userExists() {
		console.log("Checking if user exists");
		(async () => {
			try {
				const response = await fetch("/api/user-exists", {
					method: "GET",
					credentials: "include",
				});

				if (response.status === 410) {
					console.log("User does not exist");
					navigate("/signup");
					return;
				} else if (!response.ok) {
					throw new Error("Failed to check if user exists: " + response.status);
				}
				console.log("User exists");
			} catch (error) {
				console.error("Error checking if user exists:", error);
			}
		})();
	}

	useEffect(() => {
		userExists();
	}, []);

	function onSubmit(password: string) {
		console.log("Submitted");
		(async () => {
			try {
				const response = await fetch("/api/user-sign-in", {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify(password),
					credentials: "include",
				});

				if (!response.ok) {
					throw new Error("Failed to log in user: " + response.status);
				}
				console.log("Successfully logged in user");
				navigate("/dashboard");
			} catch (error) {
				console.error("Error logging in user:", error);
				setError("Invalid password");
			}
		})();
	}

	return (
		<PasswordScreen buttonText="Sign in (uses cookies)" passwordSubmit={onSubmit} error={error} />
	);
};

export default SignInScreen;
