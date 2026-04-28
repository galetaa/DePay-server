export type Role = "user" | "merchant" | "compliance" | "admin";

export type Persona = {
  id: string;
  label: string;
  subtitle: string;
  roles: Role[];
  defaultPath: string;
};

export const personas: Persona[] = [
  {
    id: "admin",
    label: "Platform Admin",
    subtitle: "Full console access",
    roles: ["admin", "merchant", "compliance", "user"],
    defaultPath: "/admin/system-health"
  },
  {
    id: "merchant",
    label: "North Coffee",
    subtitle: "Merchant operations",
    roles: ["merchant"],
    defaultPath: "/merchant/dashboard"
  },
  {
    id: "user",
    label: "Alice Stone",
    subtitle: "Wallet owner",
    roles: ["user"],
    defaultPath: "/user/dashboard"
  },
  {
    id: "compliance",
    label: "Compliance Operator",
    subtitle: "KYC and risk review",
    roles: ["compliance"],
    defaultPath: "/compliance/kyc"
  }
];

export function personaById(id: string | null | undefined) {
  return personas.find((persona) => persona.id === id) ?? personas[0];
}
