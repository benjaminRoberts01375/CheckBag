import "../styles.css";
import NavbarStyles from "./navbar.module.css";
import { useNavigate, useLocation } from "react-router-dom";
import { useList } from "../context-hook";
import Timescale from "./timescale";
import { Checkbox } from "@mui/material";
import { useState, useEffect } from "react";
import { Fade as Hamburger } from "hamburger-react";

interface NavbarProps {
	isMobileView: boolean;
}

const Navbar = ({ isMobileView }: NavbarProps) => {
	const navigate = useNavigate();
	const location = useLocation();
	const { services, serviceToggle } = useList();
	const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

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
				</div>
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
		</>
	);
};

export default Navbar;
