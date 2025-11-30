import {Client} from "@core/models/client";

export function getFullName(client: Client): string {
  return `${client.last_name} ${client.first_name} ${client.middle_name}`;
}
