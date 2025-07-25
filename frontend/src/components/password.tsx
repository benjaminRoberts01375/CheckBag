import logo from "../assets/CheckBag.svg";
import PasswordStyles from "./password.module.css";
import { FormEvent, useState, useRef, useEffect } from "react";
import AnimatedBackground from "./animated-background";

interface IPAddressProps {
	duration: number;
	waitTime: number;
}

const AnimatedIPAddress = ({ duration, waitTime }: IPAddressProps) => {
	const [ip, setIp] = useState([192, 168, 1, 1]);
	const [port, setPort] = useState(8080);
	const intervalRef = useRef<NodeJS.Timeout | null>(null);
	const animationRef = useRef<number | null>(null);

	function generateRandomIP(): number[] {
		return [
			Math.floor(Math.random() * 256),
			Math.floor(Math.random() * 256),
			Math.floor(Math.random() * 256),
			Math.floor(Math.random() * 256),
		];
	}

	function generateRandomPort(): number {
		// Generate ports in common ranges: 80-10000
		return Math.floor(Math.random() * 9920) + 80;
	}

	function animateToNewIP(newIP: number[], newPort: number) {
		// Cancel any existing animation
		if (animationRef.current) {
			cancelAnimationFrame(animationRef.current);
		}

		const startIP = [...ip];
		const startPort = port;
		const startTime = Date.now();

		const animate = () => {
			const elapsed = Date.now() - startTime;
			const progress = Math.min(elapsed / duration, 1);

			// Easing function for smooth animation
			const easeOutCubic = (t: number) => 1 - Math.pow(1 - t, 3);
			const easedProgress = easeOutCubic(progress);

			const currentIP = startIP.map((start, index) => {
				const end = newIP[index];
				const current = Math.round(start + (end - start) * easedProgress);
				return current;
			});

			const currentPort = Math.round(startPort + (newPort - startPort) * easedProgress);

			setIp(currentIP);
			setPort(currentPort);

			if (progress < 1) {
				animationRef.current = requestAnimationFrame(animate);
			} else {
				animationRef.current = null;
			}
		};

		animationRef.current = requestAnimationFrame(animate);
	}

	useEffect(() => {
		// Start the first animation immediately
		animateToNewIP(generateRandomIP(), generateRandomPort());

		// Then start the interval for subsequent animations
		intervalRef.current = setInterval(() => {
			animateToNewIP(generateRandomIP(), generateRandomPort());
		}, waitTime);

		// Cleanup on unmount
		return () => {
			if (intervalRef.current) {
				clearInterval(intervalRef.current);
			}
			if (animationRef.current) {
				cancelAnimationFrame(animationRef.current);
			}
		};
	}, []);

	return (
		<div>
			<h1 id={PasswordStyles["IPAddress"]}>
				{ip.join(".")}:{port}
			</h1>
		</div>
	);
};

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
						<div id={PasswordStyles["ipOverlay"]}>
							<AnimatedIPAddress duration={1000} waitTime={5000} />
						</div>
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
