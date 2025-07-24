import { useState, useEffect, useRef } from "react";
import PasswordStyles from "./password.module.css";
import { FormEvent } from "react";

const AnimatedBackground = () => {
	const svgRef = useRef(null);
	const [dimensions, setDimensions] = useState({
		width: window.innerWidth,
		height: window.innerHeight,
	});

	useEffect(() => {
		const handleResize = () => {
			setDimensions({ width: window.innerWidth, height: window.innerHeight });
		};

		window.addEventListener("resize", handleResize);
		return () => window.removeEventListener("resize", handleResize);
	}, []);

	// Generate random curved path that goes through center
	const generatePath = (startY: number, endY: number, width: number, height: number): string => {
		const startX = -50; // Start off-screen left
		const endX = width + 50; // End off-screen right

		// Define center area (roughly where login form is)
		const centerY = height / 2;
		const centerRadius = Math.min(width, height) * 0.2; // 20% of screen size

		// First control point - guides path toward center from start
		const cp1X = width * 0.3 + (Math.random() - 0.5) * width * 0.1;
		const cp1Y = centerY + (Math.random() - 0.5) * centerRadius;

		// Second control point - guides path away from center toward end
		const cp2X = width * 0.7 + (Math.random() - 0.5) * width * 0.1;
		const cp2Y = centerY + (Math.random() - 0.5) * centerRadius;

		return `M ${startX} ${startY} C ${cp1X} ${cp1Y}, ${cp2X} ${cp2Y}, ${endX} ${endY}`;
	};

	// Generate multiple paths
	const paths = Array.from({ length: 10 }, (_, i) => {
		const startY = Math.random() * dimensions.height;
		const endY = Math.random() * dimensions.height;
		return {
			id: i,
			path: generatePath(startY, endY, dimensions.width, dimensions.height),
			duration: 8 + Math.random() * 4, // 8-12 seconds
			delay: Math.random() * 5, // 0-5 seconds delay
		};
	});

	return (
		<div className={PasswordStyles.backgroundContainer}>
			<svg
				ref={svgRef}
				width={dimensions.width}
				height={dimensions.height}
				className={PasswordStyles.backgroundSvg}
			>
				<defs>
					{/* Glow effect for balls */}
					<filter id="glow" x="-50%" y="-50%" width="200%" height="200%">
						<feGaussianBlur stdDeviation="4" result="coloredBlur" />
						<feMerge>
							<feMergeNode in="coloredBlur" />
							<feMergeNode in="SourceGraphic" />
						</feMerge>
					</filter>
				</defs>

				{paths.map(pathData => (
					<g key={pathData.id}>
						{/* Invisible path for animation */}
						<path id={`path-${pathData.id}`} d={pathData.path} fill="none" stroke="none" />

						{/* Optional: visible path for debugging */}
						<path
							d={pathData.path}
							fill="none"
							stroke="rgba(255, 255, 255, 0.1)"
							strokeWidth="1"
							strokeDasharray="5,5"
						/>

						{/* Animated ball */}
						<circle r="12" fill="#ffff00" filter="url(#glow)" opacity="0">
							<animateMotion
								dur={`${pathData.duration}s`}
								begin={`${pathData.delay}s`}
								repeatCount="indefinite"
								rotate="auto"
							>
								<mpath href={`#path-${pathData.id}`} />
							</animateMotion>

							{/* Instant visibility animation */}
							<animate
								attributeName="opacity"
								values="1"
								dur={`${pathData.duration}s`}
								begin={`${pathData.delay}s`}
								repeatCount="indefinite"
							/>
						</circle>
					</g>
				))}
			</svg>
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
			<AnimatedBackground />
			<div id={PasswordStyles["container"]}>
				<div id={PasswordStyles["wrapper"]}>
					<div id={PasswordStyles["logo"]}>
						<div id={PasswordStyles["placeholder"]}>
							<h1>CheckBag Logo Placeholder</h1>
						</div>
					</div>
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
