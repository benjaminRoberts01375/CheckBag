import "../styles.css";
import DashboardStyles from "./dashboard.module.css";
import APIStyles from "./api.module.css";
import ServicesStyles from "./services.module.css";
import { useList } from "../context-hook";
import APIKey from "../types/api-key.tsx";
import { useState, useEffect } from "react";

const APIScreen = () => {
	const { apiKeys, requestServiceData } = useList();

	useEffect(() => {
		requestServiceData();
	}, []);

	return (
		<div id={DashboardStyles["container"]}>
			<title>CheckBag - API Keys</title>
			<div className={DashboardStyles["graph-group"]}>
				<h2 className="header">API Keys</h2>
				{apiKeys ? apiKeys.map(key => <APIKeyItem apiKey={key} key={key.id} />) : null}
				<APIKeyAdd key={"new"} />
			</div>
		</div>
	);
};

export default APIScreen;

interface APIKeyItemProps {
	apiKey: APIKey;
}

const APIKeyItem = ({ apiKey }: APIKeyItemProps) => {
	const { removeAPIKey } = useList();

	function deleteKey(): void {
		if (apiKey.id) {
			removeAPIKey(apiKey.id);
		}
	}

	return (
		<div className={APIStyles["row"]}>
			<p id={APIStyles["delete"]}>{apiKey.key ? apiKey.name + ": " : apiKey.name}</p>
			<p className={APIStyles["raw-key"]}>{apiKey.key ? apiKey.key : ""}</p>
			<button
				className={`${ServicesStyles.delete} primary ${APIStyles["delete"]}`}
				onClick={() => deleteKey()}
			>
				Delete
			</button>
		</div>
	);
};

const APIKeyAdd = () => {
	const [name, setName] = useState<string>("");
	const { addAPIKey } = useList();

	return (
		<div className={APIStyles["row"]}>
			<input
				type="text"
				placeholder="Name"
				value={name}
				onChange={e => setName(e.target.value)}
				className={`${ServiceEditStyles["input"]} ${APIStyles["text-field"]}`}
			/>
			<button
				className={`${ServiceEditStyles.submit} primary ${APIStyles["submit"]}`}
				onClick={() => {
					addAPIKey(name);
					setName("");
				}}
				disabled={name === ""}
			>
				Add
			</button>
		</div>
	);
};
