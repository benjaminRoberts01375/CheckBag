import "../styles.css";
import NavbarStyles from "./navbar.module.css";

const Navbar = () => {
	return (
		<nav id={NavbarStyles["navbar-container"]}>
			<h1>Navbar</h1>
		</nav>
	);
};

export default Navbar;
