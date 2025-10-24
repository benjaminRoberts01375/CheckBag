import "../styles.css";
import PasswordScreen from "../components/password";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useList } from "../context-hook";

const SignInScreen = () => {
	const [error, setError] = useState<string>("");
	const navigate = useNavigate();
	const { signIn } = useList();

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
				signIn();
				navigate("/dashboard");
			} catch (error) {
				console.error("Error logging in user:", error);
				setError("Invalid password");
			}
		})();
	}

	return (
		<>
			<title>CheckBag - Sign Up</title>
			<PasswordScreen buttonText="Sign in (uses cookies)" passwordSubmit={onSubmit} error={error} />
		</>
	);
};

export default SignInScreen;
