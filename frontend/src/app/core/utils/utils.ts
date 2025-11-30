import {Client} from "@core/models/client";

export function getFullName(client: Client): string {
  if (client.middle_name) return `${client.last_name} ${client.first_name} ${client.middle_name}`;
  return `${client.last_name} ${client.first_name}`
}
