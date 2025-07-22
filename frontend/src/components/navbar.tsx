import "../styles.css";
import NavbarStyles from "./navbar.module.css";
import { useNavigate, useLocation } from "react-router-dom";
import { useList } from "../context-hook";
import Timescale from "./timescale";
import { Checkbox } from "@mui/material";

const Navbar = () => {
	const navigate = useNavigate();
	const location = useLocation();
	const { services, serviceToggle } = useList();

	return (
		<div id={NavbarStyles["navbar-container"]}>
			<h1 id={NavbarStyles["title"]}>CheckBag</h1>
			<div id={NavbarStyles["divider"]}></div>
			<button
				className={`${NavbarStyles["entry"]} ${location.pathname === "/dashboard/home" ? NavbarStyles["active"] : ""}`}
				onClick={() => navigate("/dashboard/home")}
			>
				<p>Home</p>
			</button>
			<button
				className={`${NavbarStyles["entry"]} ${location.pathname === "/dashboard/services" ? NavbarStyles["active"] : ""}`}
				onClick={() => navigate("/dashboard/services")}
			>
				<p>Services</p>
			</button>
			<div id={NavbarStyles["divider"]}></div>
			<div id={NavbarStyles["services"]}>
				{services.map(service => (
					<button
						className={`${NavbarStyles["entry"]}`}
						onClick={() => serviceToggle(service.clientID)}
						key={service.clientID}
					>
						<div className={NavbarStyles["entryService"]}>
							<Checkbox
								checked={service.enabled || service.enabled === undefined}
								readOnly={true}
								sx={{
									color: "#ffd20a",
									"&.Mui-checked": {
										color: "#ffd20a",
									},
									"& .MuiSvgIcon-root": {
										fontSize: "16pt",
										margin: 0,
										padding: 0,
									},
									"&.MuiCheckbox-root": {
										padding: 0,
										margin: 0,
									},
								}}
							/>
							<p>{service.title}</p>
						</div>
					</button>
				))}
			</div>
			<Timescale />
		</div>
	);
};

export default Navbar;
