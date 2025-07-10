import "../styles.css";
import PasswordScreen from "../components/password";

const SignUpScreen = () => {
	function onSubmit(event: React.FormEvent<HTMLFormElement>) {
		event.preventDefault();
		console.log("Submitted");
	}

	return <PasswordScreen buttonText="Sign in (uses cookies)" onSubmit={onSubmit} />;
};

export default SignUpScreen;
