import "../styles.css";
import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import PasswordScreen from "../components/password";

const SignUpScreen = () => {
	const navigate = useNavigate();
	const [error, setError] = useState<string>("");

	function userExists() {
		console.log("Checking if user exists");
		(async () => {
			try {
				const response = await fetch("/api/user-exists", {
					method: "GET",
					headers: {
						"Content-Type": "application/json",
					},
					credentials: "include",
				});

				if (response.status === 410) {
					console.log("User does not exist");
					return;
				} else if (!response.ok) {
					throw new Error("Failed to check if user exists: " + response.status);
				}
				console.log("User exists");
				navigate("/signin");
			} catch (error) {
				const err = error as Error;
				console.error("Error checking if user exists: ", err.message);
				setError("Unable to check if user exists: " + err.message);
			}
		})();
	}

	useEffect(() => {
		userExists();
	}, []);

	function onSubmit(password: string) {
		console.log("Submitting user");
		(async () => {
			try {
				const response = await fetch("/api/user-sign-up", {
					method: "POST",
					body: JSON.stringify(password),
					credentials: "include",
				});

				if (!response.ok) {
					throw new Error("Failed to sign up user: " + response.status);
				}
				console.log("Successfully signed up user");
				navigate("/dashboard");
			} catch (error) {
				const err = error as Error;
				console.error("Error signing up user: ", err.message);
				setError("Unable to sign up user: " + err.message);
			}
		})();
	}

	return <PasswordScreen buttonText="Create account" passwordSubmit={onSubmit} error={error} />;
};

export default SignUpScreen;
