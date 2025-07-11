import "../styles.css";
import PasswordStyles from "./password.module.css";
import { FormEvent } from "react";

interface PasswordScreenProps {
	buttonText: string;
	passwordSubmit: (password: string) => void;
	error: string;
}

const PasswordScreen = ({ buttonText, passwordSubmit, error }: PasswordScreenProps) => {
	function onSubmit(event: FormEvent<HTMLFormElement>) {
		event.preventDefault();
		const formData = new FormData(event.currentTarget);
		passwordSubmit(formData.get("password") as string);
	}

	return (
		<div id={PasswordStyles["container"]}>
			<div id={PasswordStyles["wrapper"]}>
				<div id={PasswordStyles["logo"]}>
					<div id={PasswordStyles["placeholder"]}>
						<h1>CheckBag Logo Placeholder</h1>
					</div>
				</div>
				<form onSubmit={onSubmit}>
					<input placeholder="Password" type="password" id={PasswordStyles["field"]} />
					<button type="submit" id={PasswordStyles["submit"]}>
						{buttonText}
					</button>
				</form>
				<p>{error}</p>
			</div>
		</div>
	);
};

export default PasswordScreen;
