import "../styles.css";
import EmailStyles from "./checkEmail.module.css";
import { useLocation } from "react-router-dom";

interface LocationState {
	userEmail?: string;
}

const CheckEmail: React.FC = () => {
	const location = useLocation();
	const state = location.state as LocationState | null;
	const email = state?.userEmail;

	return (
		<p id={EmailStyles["explanation"]}>
			Check your email
			{email ? <span id={EmailStyles["bold"]}> {email}</span> : ""} for a verification link.
		</p>
	);
};

export default CheckEmail;
