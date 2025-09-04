import "../styles.css";
import ServicesStyles from "./services.module.css";
import DashboardStyles from "./dashboard.module.css";
import { useList } from "../context-hook";
import { useEffect, useState } from "react";
import Service from "../types/service.tsx";

const ServicesScreen = () => {
	const { services, requestServiceData } = useList();
	useEffect(() => {
		requestServiceData();
	}, []);

	return (
		<div id={DashboardStyles["container"]}>
			<title>CheckBag - Services</title>
			<div className={DashboardStyles["graph-group"]}>
				<h2 className="header">Services</h2>
				{services.map(service => (
					<ServiceEdit service={service} key={service.clientID} />
				))}
				<ServiceEdit service={undefined} key={"new"} />
			</div>
		</div>
	);
};

export default ServicesScreen;

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
	const { serviceAdd, serviceDelete, serviceUpdate } = useList();
	function createService() {
		serviceAdd(new Service(internalAddress, externalAddress, name));
		setName("");
		setInternalAddress("");
		setExternalAddress([""]);
	}

	function updateService(service: Service): void {
		service.internal_address = internalAddress;
		service.external_address = externalAddress;
		service.title = name;
		serviceUpdate(service);
	}

	return (
		<div id={ServicesStyles["input-container"]}>
			<input
				type="text"
				placeholder="Name"
				value={name}
				onChange={e => setName(e.target.value)}
				className={ServicesStyles["input"]}
			/>

			<input
				type="text"
				placeholder="External Address"
				value={externalAddress}
				onChange={e => setExternalAddress([e.target.value])}
				className={ServicesStyles["input"]}
			/>
			<input
				type="text"
				placeholder="Internal IP Address"
				value={internalAddress}
				onChange={e => setInternalAddress(e.target.value)}
				className={ServicesStyles["input"]}
			/>
			<div id={ServicesStyles["buttons"]}>
				{service ? (
					<>
						<button
							onClick={() => updateService(service)}
							title={"ClientID: " + service.clientID + ", ID: " + service.id}
							className={`${ServicesStyles.submit} primary`}
							disabled={
								name === service.title &&
								internalAddress === service.internal_address &&
								externalAddress[0] === service.external_address[0]
							}
						>
							Update
						</button>
						<button
							className={`${ServicesStyles.delete} primary`}
							onClick={() => serviceDelete(service.clientID)}
						>
							Delete
						</button>
					</>
				) : (
					<button
						className={`${ServicesStyles.submit} primary`}
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
