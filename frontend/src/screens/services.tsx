import "../styles.css";
import ServicesStyles from "./services.module.css";
import DashboardStyles from "./dashboard.module.css";
import { useList } from "../context-hook";
import { useState, useEffect, useRef } from "react";
import Service from "../types/service.tsx";
import { CgMoreVerticalAlt } from "react-icons/cg";

const ServicesScreen = () => {
	const { services, requestServiceData } = useList();
	useEffect(() => {
		requestServiceData();
	}, []);
	const dialogRef = useRef<HTMLDialogElement | null>(null);

	return (
		<div id={DashboardStyles["container"]}>
			<title>CheckBag - Services</title>
			<dialog ref={dialogRef}>
				<EditService service={undefined} finish={() => dialogRef.current?.close()} />
			</dialog>
			<div className={DashboardStyles["graph-group"]}>
				<h2 className="header">Services</h2>
				{services.map(service => (
					<ServiceEntry service={service} key={service.clientID} />
				))}
				<button
					className="submit"
					id={ServicesStyles["add-service-button"]}
					onClick={() => dialogRef.current?.showModal()}
				>
					Add Service
				</button>
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

	return (
		<div id={ServicesStyles["service-container"]}>
			<dialog ref={dialogRef}>
				<EditService service={service} finish={() => dialogRef.current?.close()} />
			</dialog>
			<h2>{service.title}</h2>
			<div className={ServicesStyles["connection-info"]}>
				<div id={ServicesStyles["service-endpoints"]}>
					{service.external_address.map(externalAddress => (
						<ServiceStatus address={externalAddress} key={service.clientID} />
					))}
				</div>
				{service.internal_address ? <ServiceStatus address={service.internal_address} /> : null}
				<button onClick={() => dialogRef.current?.showModal()}>
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

interface EditServiceProps {
	service: Service | undefined;
	finish: () => void;
}

const EditService = ({ service, finish }: EditServiceProps) => {
	const { serviceAdd, serviceDelete, serviceUpdate } = useList();
	const [title, setTitle] = useState(service?.title ?? "");
	const [incomingAddresses, setIncomingAddress] = useState(service?.external_address ?? []);
	const [outgoingAddress, setOutgoingAddress] = useState(service?.internal_address ?? "");

	function submit(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
		e.preventDefault();
		service == undefined ? createService() : updateService(service);
		finish();
	}

	function createService() {
		serviceAdd(new Service(outgoingAddress, incomingAddresses, title));
	}

	function updateService(service: Service) {
		service.title = title;
		service.external_address = incomingAddresses;
		service.internal_address = outgoingAddress;
		serviceUpdate(service);
	}

	function cancel(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
		e.preventDefault();
		console.log("Cancelling");
		finish();
		setTitle(service?.title ?? "");
		setIncomingAddress(service?.external_address ?? []);
		setOutgoingAddress(service?.internal_address ?? "");
	}

	function deleteService(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
		e.preventDefault();
		console.log("Deleting");
		if (service) {
			serviceDelete(service.clientID);
		}
		finish();
	}

	return (
		<form id={ServicesStyles["edit-service-container"]}>
			<div id={ServicesStyles["edit-service-header"]}>
				<h1>Editing</h1>
				<input
					type="text"
					placeholder="Service Name"
					value={title}
					onChange={e => setTitle(e.target.value)}
					className={ServicesStyles["input"]}
				/>
			</div>
			<h3>From:</h3>
			<input
				type="url"
				autoComplete="off"
				placeholder="Forward Address"
				value={incomingAddresses}
				onChange={e => setIncomingAddress([e.target.value])}
				className={ServicesStyles["input"]}
			/>
			<h3>To:</h3>
			<input
				type="url"
				autoComplete="off"
				placeholder="Forward Address"
				value={outgoingAddress}
				onChange={e => setOutgoingAddress(e.target.value)}
				className={ServicesStyles["input"]}
			/>
			<div id={ServicesStyles["submission-buttons"]}>
				{service ? (
					<button
						className="delete"
						role="delete"
						onClick={e => {
							deleteService(e);
						}}
					>
						Delete
					</button>
				) : null}

				<button
					className="cancel"
					role="cancel"
					onClick={e => {
						cancel(e);
					}}
				>
					Cancel
				</button>
				<button
					className="submit"
					role="submit"
					onClick={e => {
						submit(e);
					}}
				>
					{service == undefined ? "Create" : "Save"}
				</button>
			</div>
		</form>
	);
};
