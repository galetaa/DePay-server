import { ShieldCheck } from "lucide-react";
import { FormEvent, useState } from "react";
import { personaById, personas, type Persona } from "../auth/personas";

type LoginPageProps = {
  persona: Persona;
  onPersonaChange: (persona: Persona) => void;
  onComplete: (path: string) => void;
};

export function LoginPage({ persona, onPersonaChange, onComplete }: LoginPageProps) {
  const [selectedPersonaId, setSelectedPersonaId] = useState(persona.id);
  const selectedPersona = personaById(selectedPersonaId);

  function onSubmit(event: FormEvent) {
    event.preventDefault();
    onPersonaChange(selectedPersona);
    onComplete(selectedPersona.defaultPath);
  }

  return (
    <section className="login-page">
      <form className="login-panel" onSubmit={onSubmit}>
        <div className="login-heading">
          <span className="brand-mark">D</span>
          <div>
            <h1>DePay Console</h1>
            <p>Demo access</p>
          </div>
        </div>
        <div className="persona-grid">
          {personas.map((item) => (
            <label className={selectedPersonaId === item.id ? "persona-card active" : "persona-card"} key={item.id}>
              <input type="radio" name="persona" value={item.id} checked={selectedPersonaId === item.id} onChange={(event) => setSelectedPersonaId(event.target.value)} />
              <ShieldCheck size={18} />
              <span>
                <strong>{item.label}</strong>
                <small>{item.subtitle}</small>
              </span>
            </label>
          ))}
        </div>
        <button type="submit" className="primary-button">
          Continue
        </button>
      </form>
    </section>
  );
}
