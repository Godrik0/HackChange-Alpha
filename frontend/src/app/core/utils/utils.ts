import {Client} from "@core/models/client";
import {never} from "rxjs";

export function getFullName(client: Client) {
  return `${client.lastName} ${client.firstName} ${client.middleName}`;
}
