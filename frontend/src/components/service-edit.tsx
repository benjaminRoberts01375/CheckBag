import "../styles.css";
import ServiceEditStyles from "./service-edit.module.css";
import Service from "../types/service.tsx";
import { useState } from "react";
import { useList } from "../context-hook";

interface ServiceAddScreenProps {
	service: Service | undefined;
}

const ServiceEdit = ({ service }: ServiceAddScreenProps) => {
	const [name, setName] = useState<string>(service ? service.title : "");
	const [internalAddress, setInternalAddress] = useState<string>(
		service ? service.internal_address : "",
	);
	const [externalAddress, setExternalAddress] = useState<string[]>(
		service ? service.external_address : [""],
	);
	const { serviceAdd } = useList();
	function createService() {
		console.log("Creating service");
		serviceAdd(new Service(internalAddress, externalAddress, name));
		setName("");
		setInternalAddress("");
		setExternalAddress([""]);
	}

	function deleteService() {
		console.log("Deleting service");
	}

	function updateService() {
		console.log("Updating service");
	}

	return (
		<div id={ServiceEditStyles["input-container"]}>
			<input
				type="text"
				placeholder="Name"
				value={name}
				onChange={e => setName(e.target.value)}
				className={ServiceEditStyles["input"]}
			/>

			<input
				type="text"
				placeholder="External Address"
				value={externalAddress}
				onChange={e => setExternalAddress([e.target.value])}
				className={ServiceEditStyles["input"]}
			/>
			<input
				type="text"
				placeholder="Internal IP Address"
				value={internalAddress}
				onChange={e => setInternalAddress(e.target.value)}
				className={ServiceEditStyles["input"]}
			/>
			<div id={ServiceEditStyles["buttons"]}>
				{service ? (
					<>
						<button
							onClick={() => updateService()}
							title={"ClientID: " + service.clientID + ", ID: " + service.id}
							className={`${ServiceEditStyles.submit} primary`}
							disabled={
								name === service.title &&
								internalAddress === service.internal_address &&
								externalAddress[0] === service.external_address[0]
							}
						>
							Update
						</button>
						<button
							className={`${ServiceEditStyles.delete} primary`}
							onClick={() => deleteService()}
						>
							Delete
						</button>
					</>
				) : (
					<button
						className={`${ServiceEditStyles.submit} primary`}
						onClick={() => createService()}
						disabled={name === "" || internalAddress === "" || externalAddress[0] === ""}
					>
						Add
					</button>
				)}
			</div>
		</div>
	);
};

export default ServiceEdit;
