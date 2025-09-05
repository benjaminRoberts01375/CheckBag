import "../styles.css";
import ServicesStyles from "./services.module.css";
import DashboardStyles from "./dashboard.module.css";
import { useList } from "../context-hook";
import { useEffect, useRef } from "react";
import Service from "../types/service.tsx";
import { CgMoreVerticalAlt } from "react-icons/cg";

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
					<ServiceEntry service={service} key={service.clientID} />
				))}
				<button className={`${ServicesStyles.submit} primary`}>Add Service</button>
			</div>
		</div>
	);
};

export default ServicesScreen;

interface ServiceListEntryProps {
	service: Service;
}

const ServiceEntry = ({ service }: ServiceListEntryProps) => {
	const dialogRef = useRef<HTMLDialogElement | null>(null);

	function openDialog(): void {
		console.log("Opening dialog");
		dialogRef.current?.showModal();
	}

	function closeDialog(): void {
		dialogRef.current?.close();
	}

	return (
		<div id={ServicesStyles["service-container"]}>
			<dialog ref={dialogRef}>
				<EditService service={service} />
			</dialog>
			<h2>{service ? service.title : "Untitled Service"}</h2>
			<div className={ServicesStyles["connection-info"]}>
				<div id={ServicesStyles["service-endpoints"]}>
					{service?.external_address.map(externalAddress => (
						<ServiceStatus address={externalAddress} key={service.clientID} />
					))}
				</div>
				{service?.internal_address ? <ServiceStatus address={service.internal_address} /> : null}
				<button onClick={() => openDialog()}>
					<CgMoreVerticalAlt className="icon" id={ServicesStyles["menu-icon"]} />
				</button>
			</div>
		</div>
	);
};

interface ServiceURLProps {
	address: string;
}

const ServiceStatus = ({ address }: ServiceURLProps) => {
	return <p className={ServicesStyles["service-status"]}>{address}</p>;
};

const EditService = ({ service }: ServiceListEntryProps) => {
	return (
		<div>
			<h1>Edit Service</h1>
		</div>
	);
};
