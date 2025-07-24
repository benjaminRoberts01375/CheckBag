import { useState, useEffect, useRef } from "react";
import AnimatedBackgroundStyles from "./animated-background.module.css";

// Normalized path data (using 0-1 coordinates)
class NormalizedPathData {
	id: number;
	startY: number;
	endY: number;
	cp1X: number;
	cp1Y: number;
	cp2X: number;
	cp2Y: number;
	duration: number;
	delay: number;

	constructor(
		id: number,
		startY: number,
		endY: number,
		cp1X: number,
		cp1Y: number,
		cp2X: number,
		cp2Y: number,
		duration: number,
		delay: number,
	) {
		this.id = id;
		this.startY = startY;
		this.endY = endY;
		this.cp1X = cp1X;
		this.cp1Y = cp1Y;
		this.cp2X = cp2X;
		this.cp2Y = cp2Y;
		this.duration = duration;
		this.delay = delay;
	}

	// Convert normalized coordinates to actual path string
	toPath(width: number, height: number): string {
		const startX = -50;
		const endX = width + 50;
		const centerY = height / 2;
		const centerRadius = Math.min(width, height) * 0.2;

		const actualStartY = this.startY * height;
		const actualEndY = this.endY * height;
		const actualCp1X = this.cp1X * width;
		const actualCp1Y = centerY + (this.cp1Y - 0.5) * centerRadius;
		const actualCp2X = this.cp2X * width;
		const actualCp2Y = centerY + (this.cp2Y - 0.5) * centerRadius;

		return `M ${startX} ${actualStartY} C ${actualCp1X} ${actualCp1Y}, ${actualCp2X} ${actualCp2Y}, ${endX} ${actualEndY}`;
	}
}

const AnimatedBackground = () => {
	const [normalizedPaths, setNormalizedPaths] = useState<NormalizedPathData[]>([]);
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

	// Only generate paths once on mount
	useEffect(() => {
		if (normalizedPaths.length === 0) {
			generateNormalizedPaths();
		}
	}, [normalizedPaths.length]);

	// Generate normalized paths (only called once on mount)
	function generateNormalizedPaths(): void {
		setNormalizedPaths(() => {
			return Array.from({ length: 10 }, (_, i) => {
				const startY = Math.random(); // 0-1
				const endY = Math.random(); // 0-1

				// Control points in normalized coordinates
				const cp1X = 0.3 + (Math.random() - 0.5) * 0.1; // Around 30% with some variance
				const cp1Y = Math.random(); // 0-1, will be adjusted relative to center
				const cp2X = 0.7 + (Math.random() - 0.5) * 0.1; // Around 70% with some variance
				const cp2Y = Math.random(); // 0-1, will be adjusted relative to center

				return new NormalizedPathData(
					i,
					startY,
					endY,
					cp1X,
					cp1Y,
					cp2X,
					cp2Y,
					8 + Math.random() * 4, // 8-12 seconds
					Math.random() * 5,
				);
			});
		});
	}

	return (
		<div className={AnimatedBackgroundStyles.backgroundContainer}>
			<svg
				ref={svgRef}
				width={dimensions.width}
				height={dimensions.height}
				className={AnimatedBackgroundStyles.backgroundSvg}
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

				{normalizedPaths.map(pathData => {
					const actualPath = pathData.toPath(dimensions.width, dimensions.height);
					return (
						<g key={pathData.id}>
							{/* Invisible path for animation */}
							<path id={`path-${pathData.id}`} d={actualPath} fill="none" stroke="none" />

							{/* Optional: visible path for debugging */}
							<path
								d={actualPath}
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
					);
				})}
			</svg>
		</div>
	);
};

export default AnimatedBackground;
