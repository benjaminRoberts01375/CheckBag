import "../styles.css";
import timescaleStyles from "./timescale.module.css";
import { useList } from "../context-hook";

const Timescale = () => {
	const { timescale, setTimescale } = useList();
	return (
		<div id={timescaleStyles["container"]} className="tertiary-background">
			<button
				className={`${timescale === "hour" ? `${timescaleStyles["active"]} primary` : timescaleStyles["passive"]}`}
				onClick={() => setTimescale("hour")}
			>
				H
			</button>
			<button
				className={`${timescale === "day" ? `${timescaleStyles["active"]} primary` : timescaleStyles["passive"]}`}
				onClick={() => setTimescale("day")}
			>
				D
			</button>
			<button
				className={`${timescale === "month" ? `${timescaleStyles["active"]} primary` : timescaleStyles["passive"]}`}
				onClick={() => setTimescale("month")}
			>
				M
			</button>
			<button
				className={`${timescale === "year" ? `${timescaleStyles["active"]} primary` : timescaleStyles["passive"]}`}
				onClick={() => setTimescale("year")}
			>
				Y
			</button>
		</div>
	);
};

export default Timescale;
