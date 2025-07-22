import "../styles.css";
import timescaleStyles from "./timescale.module.css";
import { useList } from "../context-hook";

const Timescale = () => {
	const { timescale, setTimescale } = useList();
	return (
		<div>
			<h1>Timescale</h1>
		</div>
	);
};

export default Timescale;
