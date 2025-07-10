import "../styles.css";
import PasswordStyles from "./password.module.css";

interface PasswordScreenProps {
	buttonText: string;
	onSubmit: (event: React.FormEvent<HTMLFormElement>) => void;
}

const PasswordScreen = ({ buttonText, onSubmit }: PasswordScreenProps) => {
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
				{/* </div> */}
			</div>
		</div>
	);
};

export default PasswordScreen;
