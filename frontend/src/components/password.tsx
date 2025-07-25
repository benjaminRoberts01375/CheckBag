import logo from "../assets/CheckBag.svg";
import PasswordStyles from "./password.module.css";
import { FormEvent } from "react";
import AnimatedBackground from "./animated-background";

interface PasswordScreenProps {
	buttonText: string;
	passwordSubmit: (password: string) => void;
	error: string;
}

const PasswordScreen = ({ buttonText, passwordSubmit, error }: PasswordScreenProps) => {
	function onSubmit(event: FormEvent<HTMLFormElement>) {
		event.preventDefault();
		const formData = new FormData(event.currentTarget);
		const password = formData.get("password") as string;
		if (password === undefined || password === "" || password === null) {
			error = "Password cannot be empty";
			return;
		}
		passwordSubmit(password);
	}

	return (
		<>
			<AnimatedBackground nodes={10} speed={0.8} />
			<div id={PasswordStyles["container"]}>
				<div id={PasswordStyles["wrapper"]}>
					<div id={PasswordStyles["logoWrapper"]}>
						<img src={logo} alt="CheckBag Logo" draggable={false} />
						<p>Know your network inside and out</p>
					</div>
					<form onSubmit={onSubmit}>
						<input
							placeholder="Enter your password"
							type="password"
							id={PasswordStyles["field"]}
							name="password"
						/>
						<button type="submit" id={PasswordStyles["submit"]} className="primary">
							<p id={PasswordStyles["submit-text"]}>{buttonText}</p>
						</button>
						{error !== "" ? <p id={PasswordStyles["error"]}>{error}</p> : null}
					</form>
				</div>
			</div>
		</>
	);
};

export default PasswordScreen;
