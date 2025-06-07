package usuarios

import (
	"controle_cartao/infrastructure/cadastros/usuarios"
	"github.com/google/uuid"
)

// IUsuario define uma interface para os metodos de acesso a camada de dados
type IUsuario interface {
	CadastrarUsuario(req *usuarios.Usuario) error
	BuscarUsuario(nome *string) (*usuarios.Usuario, error)
	AtualizarSenhaUsuario(novaSenha *string, usuarioID *uuid.UUID) error
}
