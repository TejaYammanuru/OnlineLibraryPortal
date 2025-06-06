class User < ApplicationRecord
  include Devise::JWT::RevocationStrategies::JTIMatcher
  devise :database_authenticatable, :registerable,
         :recoverable, :validatable, :jwt_authenticatable,
         jwt_revocation_strategy: self

  validates :name, presence: true

  enum role: { member: 0, librarian: 1, admin: 2 }

  after_initialize :set_default_role, if: :new_record?

  def set_default_role
    self.role ||= :member
  end
end
