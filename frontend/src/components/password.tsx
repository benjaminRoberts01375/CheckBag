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
					<img src={logo} alt="CheckBag Logo" id={PasswordStyles["logo"]} draggable={false} />
					<form onSubmit={onSubmit}>
						<input
							placeholder="Password"
							type="password"
							id={PasswordStyles["field"]}
							name="password"
						/>
						<button type="submit" id={PasswordStyles["submit"]} className="primary">
							{buttonText}
						</button>
					</form>
					<p>{error}</p>
				</div>
			</div>
		</>
	);
};

export default PasswordScreen;
