import "../styles.css";
import NavbarStyles from "./navbar.module.css";
import { useNavigate, useLocation } from "react-router-dom";
import { useList } from "../context-hook";
import Timescale from "./timescale";
import { Checkbox } from "@mui/material";
import { useEffect } from "react";
import { Fade as Hamburger } from "hamburger-react";
import { IoHome } from "react-icons/io5";
import { FaShareAlt } from "react-icons/fa";
import { PiAirplaneTiltFill } from "react-icons/pi";

interface NavbarProps {
	isMobileView: boolean;
	isMobileMenuOpen: boolean;
	setIsMobileMenuOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

const Navbar = ({ isMobileView, isMobileMenuOpen, setIsMobileMenuOpen }: NavbarProps) => {
	const navigate = useNavigate();
	const location = useLocation();
	const { services, serviceToggle } = useList();

	// Close mobile menu when route changes
	useEffect(() => {
		setIsMobileMenuOpen(false);
	}, [location.pathname]);

	// Close mobile menu when clicking outside
	useEffect(() => {
		const handleClickOutside = (event: MouseEvent) => {
			const navbar = document.getElementById(NavbarStyles["navbar-container"]);
			const hamburger = document.getElementById("hamburger-menu"); // Use the new ID for hamburger

			if (
				isMobileMenuOpen &&
				navbar &&
				hamburger &&
				!navbar.contains(event.target as Node) &&
				!hamburger.contains(event.target as Node)
			) {
				setIsMobileMenuOpen(false);
			}
		};

		if (isMobileMenuOpen) {
			document.addEventListener("mousedown", handleClickOutside);
		}

		return () => {
			document.removeEventListener("mousedown", handleClickOutside);
		};
	}, [isMobileMenuOpen]);

	return (
		<>
			{/* Mobile overlay */}
			{isMobileMenuOpen && (
				<div
					className={NavbarStyles["mobile-overlay"]}
					onClick={() => setIsMobileMenuOpen(false)}
				></div>
			)}

			{/* Navbar */}
			<div
				id={NavbarStyles["navbar-container"]}
				className={`${isMobileMenuOpen && isMobileView ? NavbarStyles["mobile-open"] : ""}`}
			>
				<div id={NavbarStyles["title-container"]}>
					{isMobileView && ( // Only render hamburger if in mobile view
						<Hamburger toggled={isMobileMenuOpen} toggle={setIsMobileMenuOpen} />
					)}
					<h1 id={NavbarStyles["title"]} className="primary">
						CheckBag
					</h1>
					<p>
						{__CHECKBAG_VERSION__ == "dev" ? (
							<p>Development Build</p>
						) : (
							<p>{__CHECKBAG_VERSION__}</p>
						)}
					</p>
				</div>
				<div id={NavbarStyles["divider"]}></div>
				<button
					className={`${NavbarStyles["entry"]} ${location.pathname === "/dashboard/home" ? NavbarStyles["active"] : ""}`}
					onClick={() => navigate("/dashboard/home")}
				>
					<IoHome className="icon" />
					<p>Home</p>
				</button>
				<button
					className={`${NavbarStyles["entry"]} ${location.pathname === "/dashboard/services" ? NavbarStyles["active"] : ""}`}
					onClick={() => navigate("/dashboard/services")}
				>
					<FaShareAlt className="icon" />
					<p>Services</p>
				</button>
				<button
					className={`${NavbarStyles["entry"]} ${location.pathname === "/dashboard/api" ? NavbarStyles["active"] : ""}`}
					onClick={() => navigate("/dashboard/api")}
				>
					<PiAirplaneTiltFill className="icon" />
					<p>API Keys</p>
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
		</>
	);
};

export default Navbar;
