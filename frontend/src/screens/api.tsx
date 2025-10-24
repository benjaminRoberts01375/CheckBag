import "../styles.css";
import DashboardStyles from "./dashboard.module.css";
import APIStyles from "./api.module.css";
import ServicesStyles from "./services.module.css";
import { useList } from "../context-hook";
import APIKey from "../types/api-key.tsx";
import { useState } from "react";

const APIScreen = () => {
	const { apiKeys } = useList();

	return (
		<div id={DashboardStyles["container"]}>
			<title>CheckBag - API Keys</title>
			{apiKeys ? apiKeys.map(key => <APIKeyItem apiKey={key} key={key.id} />) : null}
			<APIKeyAdd key={"new"} />
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
		<div className={`${APIStyles["row"]} ${DashboardStyles["graph-group"]} secondary-background`}>
			<p id={APIStyles["delete"]}>{apiKey.key ? apiKey.name + ": " : apiKey.name}</p>
			<p className={APIStyles["raw-key"]}>{apiKey.key ? apiKey.key : ""}</p>
			<button
				className={`${ServicesStyles.delete} primary ${APIStyles["delete"]} delete`}
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
		<div className={`${APIStyles["row"]} ${DashboardStyles["graph-group"]} secondary-background`}>
			<input
				type="text"
				placeholder="API Key Name"
				value={name}
				onChange={e => setName(e.target.value)}
				className={`${ServicesStyles["input"]} ${APIStyles["text-field"]}`}
			/>
			<button
				className={`${ServicesStyles.submit} primary ${APIStyles["submit"]} submit`}
				onClick={() => {
					addAPIKey(name);
					setName("");
				}}
				disabled={name === ""}
			>
				Create Key
			</button>
		</div>
	);
};
